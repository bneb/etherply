
# EtherPly Web Application

The user-facing platform for EtherPly, including the Marketing Landing Page and the Developer Console.

## Features

-   **Marketing Site**: "Postgres for Realtime" positioning.
-   **Developer Console**: Manage Organizations and Projects.
-   **Mocked Backend**: Auth and Billing are currently mocked for local development.

## Getting Started

1.  IInstall dependencies:
    ```bash
    npm install
    ```

2.  Run the development server:
    ```bash
    npm run dev
    ```

3.  Open [http://localhost:3000](http://localhost:3000) with your browser.

## Mocked Flows

-   **Authentication**: Enter *any* email in the Login screen.
-   **Billing**: Click "Upgrade" on any plan to see the active state.
-   **Projects**: Create a project to see it added to the list (persists in memory only).
