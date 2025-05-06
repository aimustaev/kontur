import { Client, Connection } from "@temporalio/client";

export const createClient = async (): Promise<Client> => {
  const connection = await Connection.connect({
    address: process.env.TEMPORAL_ADDRESS ?? "localhost:7233",
  });

  return new Client({ connection });
}; 