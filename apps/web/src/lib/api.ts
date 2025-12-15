
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface User {
    id: string;
    email: string;
    created_at?: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

export class ApiError extends Error {
    constructor(public status: number, message: string) {
        super(message);
        this.name = 'ApiError';
    }
}

export interface Project {
    id: string;
    name: string;
    region: string;
    created_at: string;
}

export interface Plan {
    id: string;
    name: string;
    price: number | string;
    limits: string;
}

export const api = {
    async fetchWithAuth<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
        let token: string | null = null;
        if (typeof window !== 'undefined') {
            token = localStorage.getItem('nmeshed_token');
        }

        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
            ...options.headers,
        };

        try {
            const res = await fetch(`${API_URL}${endpoint}`, {
                ...options,
                headers,
            });

            if (!res.ok) {
                if (res.status === 401) {
                    // Token expired or invalid
                    this.logout();
                    if (typeof window !== 'undefined') {
                        window.location.href = '/login';
                    }
                    throw new ApiError(401, 'Session expired. Please login again.');
                }
                if (res.status === 500) throw new ApiError(500, 'System manufacture defect.');
                const errorBody = await res.text().catch(() => res.statusText);
                throw new ApiError(res.status, `API Error: ${errorBody}`);
            }

            // Some endpoints might return 204 No Content
            if (res.status === 204) return {} as T;

            return res.json();
        } catch (error) {
            if (error instanceof ApiError) throw error;
            throw new Error('Connection to Core synchronization engine failed.');
        }
    },

    async login(email: string, password: string): Promise<AuthResponse> {
        try {
            const res = await fetch(`${API_URL}/v1/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });

            if (!res.ok) {
                if (res.status === 401) throw new ApiError(401, 'Credentials rejected.');
                if (res.status === 500) throw new ApiError(500, 'System manufacture defect.');
                throw new ApiError(res.status, `Server error: ${res.statusText}`);
            }

            const data = await res.json();

            // Persist token (MVP Velocity Decision)
            if (typeof window !== 'undefined') {
                localStorage.setItem('nmeshed_token', data.token);
                localStorage.setItem('nmeshed_user', JSON.stringify(data.user));
            }

            return data;
        } catch (error) {
            if (error instanceof ApiError) throw error;
            // Network failures (ECONNREFUSED) often appear as generic TypeErrors in fetch
            throw new Error('Connection to Core synchronization engine failed.');
        }
    },

    async listProjects(): Promise<Project[]> {
        return this.fetchWithAuth<Project[]>('/v1/projects');
    },

    async createProject(name: string, region: string): Promise<Project> {
        return this.fetchWithAuth<Project>('/v1/projects', {
            method: 'POST',
            body: JSON.stringify({ name, region }),
        });
    },

    async getPlans(): Promise<Plan[]> {
        return this.fetchWithAuth<Plan[]>('/v1/billing/plans');
    },

    logout() {
        if (typeof window !== 'undefined') {
            localStorage.removeItem('nmeshed_token');
            localStorage.removeItem('nmeshed_user');
        }
    },

    getUser(): User | null {
        if (typeof window !== 'undefined') {
            const stored = localStorage.getItem('nmeshed_user');
            if (stored) {
                try {
                    return JSON.parse(stored);
                } catch {
                    return null;
                }
            }
        }
        return null;
    }
};
