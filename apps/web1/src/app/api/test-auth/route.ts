import { NextRequest, NextResponse } from 'next/server';
import { prisma } from '@/lib/prisma';
import bcrypt from 'bcryptjs';

export async function POST(request: NextRequest) {
  try {
    const { email, password } = await request.json();

    if (!email || !password) {
      return NextResponse.json({ error: 'Email and password required' }, { status: 400 });
    }

    const user = await prisma.user.findUnique({
      where: { email },
    });

    if (!user || !user.password) {
      return NextResponse.json({ error: 'User not found or no password' }, { status: 404 });
    }

    const isPasswordValid = await bcrypt.compare(password, user.password);

    if (!isPasswordValid) {
      return NextResponse.json({ error: 'Invalid password' }, { status: 401 });
    }

    console.log('User email verification status:', {
      email: user.email,
      emailVerified: user.emailVerified,
      emailVerifiedType: typeof user.emailVerified,
      emailVerifiedValue: user.emailVerified,
    });

    if (!user.emailVerified) {
      return NextResponse.json({ error: 'Please verify your email first' }, { status: 400 });
    }

    return NextResponse.json({
      success: true,
      user: {
        id: user.id,
        email: user.email,
        name: user.name,
        emailVerified: user.emailVerified,
      },
    });
  } catch (error) {
    console.error('Test auth error:', error);
    return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
  }
}