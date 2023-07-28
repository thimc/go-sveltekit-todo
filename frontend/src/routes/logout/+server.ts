import { JWT_COOKIE } from '$env/static/private';
import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

export const GET: RequestHandler = async ({ cookies }) => {
  cookies.delete(JWT_COOKIE);
  throw redirect(303, '/login');
};
