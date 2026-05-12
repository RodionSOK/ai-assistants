import client from "./client";

export const dummyLogin = async (role) => {
  const { data } = await client.post("/dummyLogin", { role });
  return data;
};

export const register = async (email, password) => {
  const { data } = await client.post("/register", { email, password });
  return data;
};

export const login = async (email, password) => {
  const { data } = await client.post("/login", { email, password });
  return data;
};