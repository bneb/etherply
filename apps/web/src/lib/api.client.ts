export type Project = {
    id: string;
    name: string;
    region: string;
    created_at: string;
};

export type Plan = {
    id: string;
    name: string;
    price: number | string;
    limits: string;
};

class ApiClient {
    private baseUrl: string;

    constructor() {
        // In real app, this comes from ENV
        this.baseUrl = typeof window !== 'undefined' ? window.location.origin.replace('3000', '8080') : 'http://localhost:8080';
        // Hack for dev: simple proxy or direct port
    }

    private async fetch<T>(path: string, options?: RequestInit): Promise<T> {
        const url = `${this.baseUrl}${path}`;
        try {
            const res = await fetch(url, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options?.headers,
                },
            });

            if (!res.ok) {
                const errorBody = await res.text();
                throw new Error(`API Error ${res.status}: ${errorBody}`);
            }

            return res.json();
        } catch (error) {
            console.error(`Fetch failed for ${path}:`, error);
            throw error;
        }
    }

    async listProjects(): Promise<Project[]> {
        return this.fetch<Project[]>('/v1/projects');
    }

    async createProject(name: string, region: string): Promise<Project> {
        return this.fetch<Project>('/v1/projects', {
            method: 'POST',
            body: JSON.stringify({ name, region }),
        });
    }

    async getPlans(): Promise<Plan[]> {
        return this.fetch<Plan[]>('/v1/billing/plans');
    }
}

export const api = new ApiClient();
