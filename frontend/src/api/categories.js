import client from "./client";

export const getCategories = async () => {
  const { data } = await client.get("/categories");
  return data;
};

export const createCategory = async ({ name, description }) => {
  const { data } = await client.post("/categories", { name, description });
  return data;
};