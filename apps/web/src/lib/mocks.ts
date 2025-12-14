
export interface User {
    id: string;
    email: string;
    name: string;
}

export const AuthService = {
    signIn: async (email: string): Promise<User> => {
        // Simulate network delay
        await new Promise((resolve) => setTimeout(resolve, 800));

        return {
            id: 'usr_mock_123',
            email,
            name: email.split('@')[0],
        };
    },
};
