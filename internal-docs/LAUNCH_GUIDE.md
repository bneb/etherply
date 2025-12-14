# Low-Cost Launch Checklist & User Guide (Supabase Edition)

This guide provides a "fool-proof" path to launching your EtherPly stack using **Supabase** as the unified backend-as-a-service (Auth + DB), while keeping hosting cheap and scalable.

> [!IMPORTANT]
> **The Stack**:
> 1. **Frontend (`apps/web`)**: Deployed on **Vercel** (Best for Next.js).
> 2. **Backend (`etherply-sync-server`)**: Deployed on **Fly.io** (For the high-performance Go/BadgerDB sync engine).
> 3. **Auth & Database**: **Supabase** (Handles User Sign-up, Login, and Project/Billing data).
> 4. **Payments**: **Stripe** (Industry standard).

---

## Phase 1: Accounts & Prerequisites

**Goal**: Get all the "boring" sign-ups out of the way.

- [ ] **Buy a Domain Name**
    - Go to [Namecheap](https://www.namecheap.com/) or [Cloudflare](https://www.cloudflare.com/).
    - Note: Cloudflare offers free SSL and fast DNS.
- [ ] **Sign up for Supabase** (Auth & DB)
    - Go to [supabase.com](https://supabase.com/) and sign up with GitHub.
    - Create a new project named "EtherPly".
    - **Save your Database Password!** You cannot see it again.
- [ ] **Sign up for Vercel** (Frontend Hosting)
    - Go to [vercel.com](https://vercel.com/signup).
- [ ] **Sign up for Stripe** (Payments)
    - Go to [stripe.com](https://stripe.com/).
- [ ] **Sign up for Fly.io** (Sync Server Hosting)
    - Download CLI: `brew install flyctl`
    - Run `fly auth signup`.
- [ ] **Sign up for NPM / PyPI** (SDK Publishing)
    - [npmjs.com](https://www.npmjs.com/)
    - [pypi.org](https://pypi.org/)

---

## Phase 2: Supabase Setup (The Core)

**Goal**: Configure "The Backend" for your web app.

- [ ] **Get API Keys**
    - In Supabase Dashboard -> Project Settings -> API.
    - Copy `Project URL` and `anon` public key.
- [ ] **Configure Auth**
    - Authentication -> Providers.
    - Enable "Email" and "Google" (optional).
    - **Redirect URLs**: Add `http://localhost:3000/**` and `https://<your-vercel-domain>.vercel.app/**`.
- [ ] **Create Tables (SQL Editor)**
    - Go to "SQL Editor" -> New Query.
    - Run this starter schema (Human Intervention: Request specific schema if needed):
      ```sql
      -- Example: Profile table linked to Auth Users
      create table profiles (
        id uuid references auth.users not null primary key,
        email text,
        full_name text,
        created_at timestamptz default now()
      );
      -- Set up Row Level Security (RLS) so users only see their own data
      alter table profiles enable row level security;
      create policy "Users can view own profile" on profiles for select using (auth.uid() = id);
      create policy "Users can update own profile" on profiles for update using (auth.uid() = id);
      ```

---

## Phase 3: Frontend Launch (Vercel)

**Goal**: Connect your Next.js app to Supabase.

- [ ] **Install Supabase in Code** (Human Intervention)
    - Run: `npm install @supabase/supabase-js @supabase/ssr`
    - Replace any Clerk code with Supabase Auth helpers.
- [ ] **Set Environment Variables on Vercel**
    - `NEXT_PUBLIC_SUPABASE_URL`: Your Project URL.
    - `NEXT_PUBLIC_SUPABASE_ANON_KEY`: Your `anon` key.
- [ ] **Deploy to Vercel**
    - Import repo -> "Deploy".

---

## Phase 4: Sync Server Launch (Fly.io)

**Goal**: Host the Go server. It needs to verify Supabase tokens.

- [ ] **JWT Verification in Go**
    - You will need to update your Go server to verify Supabase JWTs.
    - Supabase uses a specific "JWT Secret" (Found in Settings -> API).
- [ ] **Deploy to Fly.io**
    - `fly launch` (App Name: `etherply-sync`).
    - Create Volume: `fly volumes create etherply_data --size 1`.
    - **Set Secret**: `fly secrets set SUPABASE_JWT_SECRET=your_secret_here`.
    - `fly deploy`.

---

## Phase 5: Payments (Stripe)

- [ ] **Get Keys**: Publishable & Secret Key.
- [ ] **Create Product**: "Pro Plan" -> Get Price ID.
- [ ] **Supabase + Stripe Sync** (Advanced)
    - You can use Supabase Edge Functions to listen to Stripe Webhooks (e.g., `checkout.session.completed`) to update the user's "subscription_status" column in your database automatically.

---

## Phase 6: Publish SDKs

### JavaScript SDK (`packages/sdk-js`)
> [!WARNING]
> Check package name availability on NPM. Use `@your-scope/package` if needed.
- [ ] `npm login`
- [ ] `npm publish --access public`

### Python SDK (`packages/sdk-python`)
- [ ] `pip install twine build`
- [ ] `python -m build`
- [ ] `twine upload dist/*`

---

## Summary of Costs (Supabase Edition)

| Service | Tier | Cost | Notes |
| :--- | :--- | :--- | :--- |
| **Supabase** | Free | **$0** | Up to 500MB DB, 50,000 MAU. Pauses after 1 week inactivity (Free tier only). |
| **Vercel** | Hobby | **$0** | Generous free tier. |
| **Fly.io** | Hobby | **~$5** | ~$2-5/mo for persistent volume and always-on VM. |
| **Stripe** | Standard | **0%** | Pay per transaction. |

**Total Estimated Monthly Cost**: ~$5.00 USD.

### Cost Projection (Next 3 Months)

Based on your growth plan (3 -> 15 -> 20 users), your infrastructure costs will remain flat because you are well within the "Free Tier" limits of Vercel and Supabase.

| Month | Users | Estimated Fixed Cost | Breakdown |
| :--- | :--- | :--- | :--- |
| **Month 1** | 3 | **~$5.00** | $0 Vercel + $0 Supabase + ~$4 Fly.io + ~$1 Domain |
| **Month 2** | 15 | **~$5.00** | Same. 15 users is < 0.1% of Supabase Free Tier. |
| **Month 3** | 20 | **~$5.00** | Same. No scaling events triggered yet. |

> [!TIP]
> **When do costs increase?**
> - **Supabase**: If you exceed 500MB database size (unlikely with just text data) or need >50,000 monthly active users.
> - **Fly.io**: If you need more RAM/CPU for the sync server. The $2/mo VM supports ~1k concurrent connections easily.
> - **Vercel**: Free tier is very generous. You rarely pay until you are a commercial team.
