import { useState } from "react";
import {
  Card,
  Typography,
  List,
  Avatar,
  Space,
  Input,
  Button,
  Form,
} from "antd";
import { UserOutlined, SendOutlined } from "@ant-design/icons";
import axios from "axios";

const { Title, Text } = Typography;

interface Message {
  id: string;
  content?: string;
  text?: string;
  body?: string;
  sender?: string;
  created_at: string;
  is_agent: boolean;
}

interface TicketChatProps {
  ticketId: string;
  messages: Message[];
  onNewMessage: (message: Message) => void;
}

export function TicketChat({
  ticketId,
  messages,
  onNewMessage,
}: TicketChatProps) {
  const [sending, setSending] = useState(false);
  const [form] = Form.useForm();

  const handleSendMessage = async (values: { message: string }) => {
    if (!ticketId || !values.message.trim()) return;

    setSending(true);
    try {
      const response = await axios.post(`/api/ticket/${ticketId}/messages`, {
        content: values.message.trim(),
      });

      onNewMessage(response.data);
      form.resetFields();
    } catch (error) {
      console.error("Error sending message:", error);
    } finally {
      setSending(false);
    }
  };

  return (
    <Card
      style={{
        flex: 1,
        display: "flex",
        flexDirection: "column",
        height: "100%",
        minHeight: 0, // важно для flex
        position: "relative",
        padding: 0,
      }}
      bodyStyle={{
        display: "flex",
        flexDirection: "column",
        flex: 1,
        minHeight: 0,
        padding: 0,
      }}
    >
      <div style={{ padding: "24px 24px 0 24px" }}>
        <Title level={4} style={{ marginBottom: 16 }}>Переписка</Title>
      </div>
      <div
        style={{
          flex: 1,
          overflowY: "auto",
          padding: "0 24px 0 24px",
          minHeight: 0,
        }}
      >
        <List
          itemLayout="horizontal"
          dataSource={messages}
          renderItem={(message) => (
            <List.Item
              style={{
                justifyContent: message.is_agent ? "flex-start" : "flex-end",
                padding: "8px 0",
              }}
            >
              <div
                style={{
                  maxWidth: "70%",
                  backgroundColor: message.is_agent ? "#f0f2f5" : "#e6f7ff",
                  padding: "12px",
                  borderRadius: "8px",
                }}
              >
                <Space>
                  <Avatar
                    icon={<UserOutlined />}
                    style={{
                      backgroundColor: message.is_agent ? "#1890ff" : "#52c41a",
                    }}
                  />
                  <div>
                    <Text strong>
                      {message.sender ||
                        (message.is_agent ? "Агент" : "Пользователь")}
                    </Text>
                    <div style={{ whiteSpace: "pre-wrap" }}>
                      {message.content || message.text || message.body}
                    </div>
                    <Text type="secondary" style={{ fontSize: "12px" }}>
                      {new Date(message.created_at).toLocaleString()}
                    </Text>
                  </div>
                </Space>
              </div>
            </List.Item>
          )}
        />
      </div>
      <Form
        form={form}
        onFinish={handleSendMessage}
        style={{
          borderTop: "1px solid #f0f0f0",
          padding: "16px 24px",
          background: "#fff",
          margin: 0,
        }}
      >
        <div style={{ display: 'flex', gap: 8 }}>
          <Form.Item
            name="message"
            style={{ flex: 1, marginBottom: 0 }}
            rules={[{ required: true, message: "Введите сообщение" }]}
          >
            <Input.TextArea
              placeholder="Введите сообщение..."
              autoSize={false}
              style={{
                height: 40,
                resize: 'none',
                borderRadius: '8px 0 0 8px',
                paddingTop: 8,
              }}
              onPressEnter={(e) => {
                if (!e.shiftKey) {
                  e.preventDefault();
                  form.submit();
                }
              }}
            />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button
              type="primary"
              icon={<SendOutlined />}
              size="large"
              loading={sending}
              htmlType="submit"
              style={{
                height: 40,
                borderRadius: '0 8px 8px 0',
                padding: '0 24px',
                display: 'flex',
                alignItems: 'center',
              }}
            >
              Отправить
            </Button>
          </Form.Item>
        </div>
      </Form>
    </Card>
  );
}
