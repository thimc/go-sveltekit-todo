import { JWT_COOKIE } from '$env/static/private';
import { redirect } from '@sveltejs/kit';

export const GET = async ({ cookies }) => {
	cookies.delete(JWT_COOKIE);
	throw redirect(303, '/login');
};
