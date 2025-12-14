# Manual Intervention: Stripe Billing Integration

**Goal:** Enable real billing for EtherPly.

## 1. Prerequisites
- Create a [Stripe Account](https://stripe.com).
- Get your **Secret Key** (`sk_test_...`) and **Publishable Key** (`pk_test_...`).

## 2. Backend (Go) Implementation
You must implement `HandleSubscribe` in `controlplane.go` to call Stripe API.

### Add Stripe SDK
```bash
go get -u github.com/stripe/stripe-go/v76
```

### Update `controlplane.go`
```go
import "github.com/stripe/stripe-go/v76"
import "github.com/stripe/stripe-go/v76/checkout/session"

func (h *Handler) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
    stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

    params := &stripe.CheckoutSessionParams{
        Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                Price: stripe.String("price_12345"), // Real Price ID from Stripe Dashboard
                Quantity: stripe.Int64(1),
            },
        },
        SuccessURL: stripe.String("http://localhost:3000/dashboard?success=true"),
        CancelURL:  stripe.String("http://localhost:3000/billing?canceled=true"),
    }

    s, err := session.New(params)
    // ... handle err ...
    
    json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
}
```

## 3. Frontend Implementation
Update `apps/web/src/app/(dashboard)/billing/page.tsx` to handle the redirect.

```tsx
<Button onClick={() => {
    const res = await api.subscribe('pro');
    window.location.href = res.url;
}}>
    Upgrade to Pro
</Button>
```

**Note:** This was deliberately omitted from the "One-Shot" execution to avoid introducing external dependencies that might break the build if keys are missing.
