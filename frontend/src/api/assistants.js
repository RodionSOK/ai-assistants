import client from "./client";

export const getAssistants = async (params = {}) => {
  const { data } = await client.get("/assistants", { params });
  return data;
};

export const getAssistant = async (id) => {
  const { data } = await client.get(`/assistants/${id}`);
  return data;
};

export const createAssistant = async (payload) => {
  const { data } = await client.post("/assistants", payload);
  return data;
};

export const updateAssistant = async (id, payload) => {
  const { data } = await client.put(`/assistants/${id}`, payload);
  return data;
};