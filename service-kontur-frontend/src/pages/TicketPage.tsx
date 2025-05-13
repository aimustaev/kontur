import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Layout, Spin, Button } from "antd";
import { ArrowLeftOutlined } from "@ant-design/icons";
import axios from "axios";
import { TicketInfo } from "../components/TicketInfo/TicketInfo";
import { TicketChat } from "../components/TicketChat/TicketChat";

const { Content } = Layout;

interface TicketDetails {
  id: string;
  status: string;
  user: string;
  agent?: string;
  channel: string;
  created_at: string;
  updated_at: string;
  problem_id?: number;
  vertical_id?: number;
  skill_id?: number;
  user_group_id?: number;
}

interface Message {
  id: string;
  content?: string;
  text?: string;
  body?: string;
  sender?: string;
  created_at: string;
  is_agent: boolean;
}

export function TicketPage() {
  const { ticketId } = useParams<{ ticketId: string }>();
  const navigate = useNavigate();
  const [ticket, setTicket] = useState<TicketDetails | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchTicketData = async () => {
      try {
        const [ticketResponse, messagesResponse] = await Promise.all([
          axios.get(`/api/ticket/${ticketId}`),
          axios.get(`/api/ticket/${ticketId}/messages`),
        ]);

        console.log("Ticket response:", ticketResponse.data);
        console.log("Messages response:", messagesResponse.data);

        setTicket(ticketResponse.data);
        const messagesData =
          messagesResponse.data.messages || messagesResponse.data || [];
        console.log("Parsed messages:", messagesData);
        setMessages(messagesData);
      } catch (error) {
        console.error("Error fetching ticket data:", error);
      } finally {
        setLoading(false);
      }
    };

    if (ticketId) {
      fetchTicketData();
    }
  }, [ticketId]);

  const handleNewMessage = (message: Message) => {
    setMessages((prev) => [...prev, message]);
  };

  if (loading) {
    return (
      <div
        style={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          height: "100vh",
        }}
      >
        <Spin size="large" />
      </div>
    );
  }

  if (!ticket) {
    return <div>Тикет не найден</div>;
  }

  return (
    <Layout style={{ minHeight: "100vh", background: "#f0f2f5" }}>
      <Layout.Header style={{ background: "#fff", padding: "0 24px", display: "flex", alignItems: "center" }}>
        <Button 
          icon={<ArrowLeftOutlined />} 
          onClick={() => navigate("/")}
          style={{ marginRight: "16px" }}
        >
          К списку тикетов
        </Button>
      </Layout.Header>
      <Content style={{ padding: "24px 24px 0 24px", display: "flex", gap: "24px", height: 'calc(100vh - 64px)', minHeight: 0, overflow: 'hidden' }}>
        <TicketInfo ticket={ticket} />
        <TicketChat
          ticketId={ticket.id}
          messages={messages}
          onNewMessage={handleNewMessage}
        />
      </Content>
    </Layout>
  );
}
