
// Stubbed User Data
export type User = {
    id: string;
    email: string;
    name: string;
    avatarUrl: string;
};

// Stubbed Organization
export type Organization = {
    id: string;
    name: string;
    slug: string;
    plan: 'free' | 'scale_up' | 'enterprise';
};

// Stubbed Project
export type Project = {
    id: string;
    name: string;
    region: string;
    activeConnections: number;
    status: 'healthy' | 'degraded';
};

// Mock Auth Service
export const AuthService = {
    signIn: async (email: string) => {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 800));

        // Always return a mock user
        return {
            id: 'usr_123',
            email,
            name: email.split('@')[0],
            avatarUrl: `https://api.dicebear.com/7.x/avataaars/svg?seed=${email}`
        };
    },

    getSession: () => {
        if (typeof window !== 'undefined') {
            const user = localStorage.getItem('etherply_user');
            return user ? JSON.parse(user) : null;
        }
        return null;
    }
};

// Mock Billing Service
export const BillingService = {
    getPlans: () => [
        { id: 'free', name: 'Hobby', price: 0 },
        { id: 'scale_up', name: 'Scale-Up', price: 499 },
        { id: 'enterprise', name: 'Enterprise', price: 'Custom' }
    ],

    subscribe: async (planId: string) => {
        await new Promise(resolve => setTimeout(resolve, 1500));
        return { status: 'active', clientSecret: 'mock_stripe_secret' };
    }
};

// Mock Project Service
export const ProjectService = {
    listProjects: async (orgId: string): Promise<Project[]> => {
        await new Promise(resolve => setTimeout(resolve, 500));
        return [
            { id: 'prj_1', name: 'Production App', region: 'us-east-1', activeConnections: 1240, status: 'healthy' },
            { id: 'prj_2', name: 'Staging Environment', region: 'eu-central-1', activeConnections: 45, status: 'healthy' },
            { id: 'prj_3', name: 'Dev Sandbox', region: 'us-west-2', activeConnections: 2, status: 'degraded' },
        ];
    },

    createProject: async (name: string, region: string) => {
        await new Promise(resolve => setTimeout(resolve, 1000));
        return {
            id: `prj_${Math.random().toString(36).substr(2, 9)}`,
            name,
            region,
            activeConnections: 0,
            status: 'healthy'
        };
    }
};
