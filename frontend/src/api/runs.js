import client from "./client";

export const runAssistant = async (assistantId, userPrompt) => {
  const { data } = await client.post(`/assistants/${assistantId}/run`, {
    userPrompt,
  });
  return data;
};

export const getMyRuns = async (params = {}) => {
  const { data } = await client.get("/runs/my", { params });
  return data;
};

export const getAllRuns = async (params = {}) => {
  const { data } = await client.get("/admin/runs", { params });
  return data;
};