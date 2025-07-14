import { NextRequest, NextResponse } from 'next/server';
import { authOptions } from '@/lib/auth';

export async function POST(request: NextRequest) {
  try {
    const { email, password } = await request.json();

    // Test the credentials provider directly
    const credentialsProvider = authOptions.providers?.find(
      (p) => p.id === 'credentials'
    ) as any;

    if (!credentialsProvider) {
      return NextResponse.json({ error: 'Credentials provider not found' }, { status: 500 });
    }

    const user = await credentialsProvider.authorize(
      { email, password },
      { query: {}, cookies: {}, headers: {} }
    );

    if (!user) {
      return NextResponse.json({ error: 'Authentication failed' }, { status: 401 });
    }

    return NextResponse.json({
      success: true,
      user,
    });
  } catch (error) {
    console.error('Test signin error:', error);
    return NextResponse.json({ 
      error: 'Authentication failed', 
      message: error instanceof Error ? error.message : 'Unknown error' 
    }, { status: 401 });
  }
}