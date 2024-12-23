import { useEffect, useState } from "react";
import { UserSchema } from "./types/user";

export default function Ping() {
  const [messages, setMessages] = useState<string[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/api/handler/ws");

    socket.onopen = () => {
      console.log("WebSocket connected");
      const authData = localStorage.getItem("user");
      if (authData) {
        const parsedUser = JSON.parse(authData);
        try {
          const user = UserSchema.parse(parsedUser);
          socket.send(`TASK/AUTHENTICATE_USER/${JSON.stringify(user)}`);
        } catch (error) {
          console.error(error);
          socket.close();
        }
      }
      setIsConnected(true);
    };

    socket.onmessage = (event) => {
      if (event.data === "ping") {
        socket.send("pong");
        return;
      }
      console.log("Message from server:", event.data);
      setMessages((prevMessages) => [...prevMessages, event.data]);
    };

    socket.onclose = () => {
      console.log("WebSocket connection closed");
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      socket.close();
    };
  }, []);
  useEffect(() => {
    console.log(messages, isConnected);
  }, [messages, isConnected]);
  return <div>Ping</div>;
}
