# Low-Cost Launch Checklist & User Guide (Supabase Edition)

This guide provides a "fool-proof" path to launching the **nMeshed** stack using **Supabase** as the unified backend-as-a-service (Auth + DB), while keeping hosting cheap and scalable.

> [!IMPORTANT]
> **The Stack**:
> 1. **Frontend (`apps/web`)**: Deployed on **Vercel** (Best for Next.js).
> 2. **Backend (`etherply-sync-server`)**: Deployed on **Fly.io** (For the high-performance Go/BadgerDB sync engine).
> 3. **Auth & Database**: **Supabase** (Handles User Sign-up, Login, and Project/Billing data).
> 4. **Payments**: **Stripe** (Industry standard).

---

## Phase 1: Accounts & Prerequisites

**Goal**: Get all the "boring" sign-ups out of the way.

- [x] **Buy a Domain Name** → `nmeshed.com` ✅
- [ ] **Sign up for Supabase** (Auth & DB)
    - Go to [supabase.com](https://supabase.com/) and sign up with GitHub.
    - Create a new project named "nMeshed".
    - **Save your Database Password!** You cannot see it again.
- [ ] **Sign up for Vercel** (Frontend Hosting)
    - Go to [vercel.com](https://vercel.com/signup).
- [ ] **Sign up for Stripe** (Payments)
    - Go to [stripe.com](https://stripe.com/).
- [ ] **Sign up for Fly.io** (Sync Server Hosting)
    - Download CLI: `brew install flyctl`
    - Run `fly auth signup`.
- [x] **Sign up for NPM** → `npmjs.com/package/nmeshed` ✅

---

## Phase 2: Supabase Setup (The Core)

**Goal**: Configure "The Backend" for your web app.

- [ ] **Get API Keys**
    - In Supabase Dashboard -> Project Settings -> API.
    - Copy `Project URL` and `anon` public key.
- [ ] **Configure Auth**
    - Authentication -> Providers.
    - Enable "Email" and "Google" (optional).
    - **Redirect URLs**: Add `http://localhost:3000/**` and `https://nmeshed.com/**`.
- [ ] **Create Tables (SQL Editor)**
    - Go to "SQL Editor" -> New Query.
    - Run starter schema (see SQL below).

```sql
-- Profile table linked to Auth Users
create table profiles (
  id uuid references auth.users not null primary key,
  email text,
  full_name text,
  created_at timestamptz default now()
);
alter table profiles enable row level security;
create policy "Users can view own profile" on profiles for select using (auth.uid() = id);
create policy "Users can update own profile" on profiles for update using (auth.uid() = id);
```

---

## Phase 3: Frontend Launch (Vercel)

**Goal**: Connect your Next.js app to Supabase.

- [ ] **Install Supabase in Code**
    - Run: `npm install @supabase/supabase-js @supabase/ssr`
- [ ] **Set Environment Variables on Vercel**
    - `NEXT_PUBLIC_SUPABASE_URL`: Your Project URL.
    - `NEXT_PUBLIC_SUPABASE_ANON_KEY`: Your `anon` key.
- [ ] **Deploy to Vercel**
    - Import repo -> "Deploy".
- [ ] **Add Custom Domain**
    - Point `nmeshed.com` to Vercel.

---

## Phase 4: Sync Server Launch (Fly.io)

**Goal**: Host the Go server.

- [ ] **Deploy to Fly.io**
    - `fly launch` (App Name: `nmeshed-sync`).
    - Create Volume: `fly volumes create nmeshed_data --size 1`.
    - **Set Secret**: `fly secrets set SUPABASE_JWT_SECRET=your_secret_here`.
    - `fly deploy`.

---

## Phase 5: Payments (Stripe)

- [ ] **Get Keys**: Publishable & Secret Key.
- [ ] **Create Product**: "Pro Plan" -> Get Price ID.

---

## Phase 6: Publish SDKs

### JavaScript SDK
- [x] Published to NPM as `nmeshed` ✅

### Python SDK (`packages/sdk-python`)
- [ ] `pip install twine build`
- [ ] `python -m build`
- [ ] `twine upload dist/*`

---

## Summary of Costs

| Service | Tier | Cost |
| :--- | :--- | :--- |
| **Supabase** | Free | **$0** |
| **Vercel** | Hobby | **$0** |
| **Fly.io** | Hobby | **~$5** |
| **Stripe** | Standard | **0%** |

**Total Estimated Monthly Cost**: ~$5.00 USD.
