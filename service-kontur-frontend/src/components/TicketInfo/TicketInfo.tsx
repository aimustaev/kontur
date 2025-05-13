import { Card, Typography, Tag, Descriptions } from 'antd';

const { Title } = Typography;

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

const statusColors: Record<string, string> = {
  'new': 'blue',
  'in_progress': 'orange',
  'resolved': 'green',
  'closed': 'gray',
};

interface TicketInfoProps {
  ticket: TicketDetails;
}

export function TicketInfo({ ticket }: TicketInfoProps) {
  return (
    <Card style={{ width: 300, height: 'fit-content' }}>
      <Title level={4}>Информация о тикете</Title>
      <Descriptions column={1} size="small">
        <Descriptions.Item label="ID">{ticket.id}</Descriptions.Item>
        <Descriptions.Item label="Статус">
          <Tag color={statusColors[ticket.status] || 'default'}>
            {ticket.status.toUpperCase()}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Пользователь">{ticket.user}</Descriptions.Item>
        <Descriptions.Item label="Агент">{ticket.agent || '-'}</Descriptions.Item>
        <Descriptions.Item label="Канал">
          <Tag color="purple">{ticket.channel.toUpperCase()}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label="Создан">
          {new Date(ticket.created_at).toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="Обновлен">
          {new Date(ticket.updated_at).toLocaleString()}
        </Descriptions.Item>
      </Descriptions>
    </Card>
  );
} 