import { NextResponse } from 'next/server';
import { prisma } from '@/lib/prisma';

export async function GET() {
  try {
    // Test database connection
    const userCount = await prisma.user.count();
    const verificationTokenCount = await prisma.verificationToken.count();
    
    return NextResponse.json({
      status: 'success',
      database: 'connected',
      users: userCount,
      verificationTokens: verificationTokenCount,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    console.error('Status check error:', error);
    return NextResponse.json(
      { 
        status: 'error', 
        message: 'Database connection failed',
        timestamp: new Date().toISOString(),
      },
      { status: 500 }
    );
  }
}