export const useApi = () => {
  const config = useRuntimeConfig();
  const baseURL = config.public.apiBase;

  const getToken = () => {
    if (import.meta.client) {
      return localStorage.getItem("picmap_token");
    }
    return null;
  };

  const request = async <T>(path: string, options: any = {}): Promise<T> => {
    const token = getToken();
    const headers: Record<string, string> = {
      ...options.headers,
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    if (!(options.body instanceof FormData) && options.body) {
      headers["Content-Type"] = "application/json";
    }

    return await $fetch<T>(`${baseURL}${path}`, {
      ...options,
      headers,
    });
  };

  return {
    get: <T>(path: string) => request<T>(path),
    post: <T>(path: string, body?: any) =>
      request<T>(path, { method: "POST", body }),
    put: <T>(path: string, body?: any) =>
      request<T>(path, { method: "PUT", body }),
    delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
    upload: <T>(path: string, formData: FormData) =>
      request<T>(path, { method: "POST", body: formData }),
  };
};
