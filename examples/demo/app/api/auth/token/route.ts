import { SignJWT } from 'jose';
import { NextResponse } from 'next/server';

export async function GET(request: Request) {
    const secret = process.env.ETHERPLY_JWT_SECRET;
    if (!secret) {
        // In a real app we would log this to Sentry/etc.
        console.error("Missing ETHERPLY_JWT_SECRET in Next.js environment");
        return NextResponse.json({ error: 'Server authentication not configured' }, { status: 500 });
    }

    const { searchParams } = new URL(request.url);
    const userId = searchParams.get('userId') || 'anon';

    try {
        const token = await new SignJWT({ sub: userId })
            .setProtectedHeader({ alg: 'HS256' })
            .setIssuedAt()
            .setExpirationTime('24h')
            .sign(new TextEncoder().encode(secret));

        return NextResponse.json({ token });
    } catch (err) {
        console.error("Token signing failed:", err);
        return NextResponse.json({ error: 'Internal Server Error' }, { status: 500 });
    }
}
