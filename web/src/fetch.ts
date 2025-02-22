const API_BASE_URL = "http://localhost:5002";

export const fetchApi = async (path: string, options: RequestInit = {}) =>
  await fetch(`${API_BASE_URL}${path}`, options);
