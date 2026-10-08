interface User {
  id: number;
  username: string;
  email: string;
  is_superuser: boolean;
  is_admin: boolean;
  is_guide: boolean;
  created_at: string;
}

interface LoginResponse {
  token: string;
  user: User;
}

export const useAuth = () => {
  const api = useApi();

  const user = useState<User | null>("auth_user", () => null);
  const token = useState<string | null>("auth_token", () => null);

  const isAuthenticated = computed(() => !!token.value);
  const isSuperuser = computed(() => user.value?.is_superuser ?? false);
  const isAdmin = computed(() => user.value?.is_admin ?? false);
  const isGuide = computed(() => user.value?.is_guide ?? false);

  const loadFromStorage = () => {
    if (import.meta.client) {
      const savedToken = localStorage.getItem("picmap_token");
      const savedUser = localStorage.getItem("picmap_user");
      if (savedToken) token.value = savedToken;
      if (savedUser) {
        try {
          user.value = JSON.parse(savedUser);
        } catch (e) {
          console.error("Failed to parse user", e);
        }
      }
    }
  };

  const login = async (username: string, password: string) => {
    const res = await api.post<LoginResponse>("/api/auth/login", {
      username,
      password,
    });

    token.value = res.token;
    user.value = res.user;

    if (import.meta.client) {
      localStorage.setItem("picmap_token", res.token);
      localStorage.setItem("picmap_user", JSON.stringify(res.user));
    }

    return res;
  };

  const logout = () => {
    token.value = null;
    user.value = null;
    if (import.meta.client) {
      localStorage.removeItem("picmap_token");
      localStorage.removeItem("picmap_user");
    }
  };

  const fetchMe = async () => {
    try {
      const me = await api.get<User>("/api/auth/me");
      user.value = me;
      if (import.meta.client) {
        localStorage.setItem("picmap_user", JSON.stringify(me));
      }
      return me;
    } catch (e) {
      logout();
      throw e;
    }
  };

  return {
    user,
    token,
    isAuthenticated,
    isSuperuser,
    isAdmin,
    isGuide,
    login,
    logout,
    fetchMe,
    loadFromStorage,
  };
};
