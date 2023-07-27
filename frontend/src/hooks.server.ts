import { authenticateUser } from '$lib/server/auth';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ resolve, event }) => {
  event.locals.user = await authenticateUser(event);
  console.log('Hook handle - Auth:', event.locals.user?.token);
  console.log('Cookie:', event.cookies.get('jwt'));

  /* Protect everything except login */
  if (event.url.pathname != '/login') {
    if (!event.locals.user) {
      throw redirect(303, '/login');
    }
  }

	return await resolve(event);
};
