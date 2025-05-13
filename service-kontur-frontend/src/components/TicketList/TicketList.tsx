import { useEffect, useState } from "react";
import { Table, Layout, Typography, Tag, Card } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const { Header, Content } = Layout;
const { Title } = Typography;

interface Ticket {
  id: string;
  status: string;
  user: string;
  agent?: string;
  problem_id?: number;
  vertical_id?: number;
  skill_id?: number;
  user_group_id?: number;
  channel: string;
  created_at: string;
  updated_at: string;
}

const statusColors: Record<string, string> = {
  new: "blue",
  in_progress: "orange",
  resolved: "green",
  closed: "gray",
};

export function TicketList() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchTickets = async () => {
      try {
        const response = await axios.get("/api/tickets");
        console.log("API Response:", response.data);

        const ticketsData = Array.isArray(response.data.tickets)
          ? response.data.tickets
          : [];
        setTickets(ticketsData);
      } catch (error) {
        console.error("Error fetching tickets:", error);
        setTickets([]);
      } finally {
        setLoading(false);
      }
    };

    fetchTickets();
  }, []);

  const columns: ColumnsType<Ticket> = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
      width: 220,
      sorter: (a, b) => a.id.localeCompare(b.id),
    },
    {
      title: "Статус",
      dataIndex: "status",
      key: "status",
      render: (status: string) => (
        <Tag color={statusColors[status] || "default"}>
          {status.toUpperCase()}
        </Tag>
      ),
      filters: Object.keys(statusColors).map((status) => ({
        text: status.toUpperCase(),
        value: status,
      })),
      onFilter: (value, record) => record.status === value,
      sorter: (a, b) => a.status.localeCompare(b.status),
    },
    {
      title: "Пользователь",
      dataIndex: "user",
      key: "user",
      sorter: (a, b) => a.user.localeCompare(b.user),
      filterSearch: true,
      filters: Array.from(new Set(tickets.map((ticket) => ticket.user))).map(
        (user) => ({
          text: user,
          value: user,
        })
      ),
      onFilter: (value, record) => record.user === value,
    },
    {
      title: "Агент",
      dataIndex: "agent",
      key: "agent",
      render: (agent?: string) => agent || "-",
      sorter: (a, b) => (a.agent || "").localeCompare(b.agent || ""),
      filters: Array.from(
        new Set(tickets.map((ticket) => ticket.agent).filter(Boolean))
      )
        .filter((agent): agent is string => agent !== undefined)
        .map((agent) => ({
          text: agent,
          value: agent,
        })),
      onFilter: (value, record) => record.agent === value,
    },
    {
      title: "Канал",
      dataIndex: "channel",
      key: "channel",
      render: (channel: string) => (
        <Tag color="purple">{channel.toUpperCase()}</Tag>
      ),
      sorter: (a, b) => a.channel.localeCompare(b.channel),
      filters: Array.from(new Set(tickets.map((ticket) => ticket.channel))).map(
        (channel) => ({
          text: channel.toUpperCase(),
          value: channel,
        })
      ),
      onFilter: (value, record) => record.channel === value,
    },
    {
      title: "Создан",
      dataIndex: "created_at",
      key: "created_at",
      render: (date: string) => new Date(date).toLocaleString(),
      sorter: (a, b) =>
        new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
    },
    {
      title: "Обновлен",
      dataIndex: "updated_at",
      key: "updated_at",
      render: (date: string) => new Date(date).toLocaleString(),
      sorter: (a, b) =>
        new Date(a.updated_at).getTime() - new Date(b.updated_at).getTime(),
    },
  ];

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Header style={{ background: "#fff", padding: "0 24px" }}>
        <Title level={3} style={{ margin: "16px 0" }}>
          Система тикетов Контур
        </Title>
      </Header>
      <Content style={{ padding: "24px" }}>
        <Card>
          <Table
            columns={columns}
            dataSource={tickets}
            rowKey="id"
            loading={loading}
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showTotal: (total) => `Всего тикетов: ${total}`,
            }}
            onChange={(pagination, filters, sorter, extra) => {
              console.log("Table parameters:", {
                pagination,
                filters,
                sorter,
                extra,
              });
            }}
            onRow={(record) => ({
              onClick: () => navigate(`/tickets/${record.id}`),
              style: { cursor: "pointer" },
            })}
          />
        </Card>
      </Content>
    </Layout>
  );
}
