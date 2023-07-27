import { API_URL } from '$env/static/private';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	const { user } = locals;

	const getTodos = async () => {
		console.log('Fetching todos');
		const token = user?.token ?? '';
		try {
			const res = await fetch(`${API_URL}/api/v1/todos`, {
				method: 'GET',
				headers: { Authorization: `Bearer ${token}` }
			});
			const result = await res.json();

			const todos = result.result;
			return todos;
		} catch (err) {
			console.log('getTodos error:', err);
		}
		return null;
	};

	return {
		todos: getTodos(),
		user: user
	};
};

export const actions: Actions = {
	createTodo: async ({ locals, request }) => {
		const formData = await request.formData();
		const content = formData.get('content');
		const user = locals.user;

		try {
			const res = await fetch(`${API_URL}/api/v1/todos`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${user?.token}`
				},
				body: JSON.stringify({
					content,
					createdBy: user?.id,
					title: 'Todo'
				})
			});
			const result = await res.json();
			console.log('Create todo result:', result);
			if (result.success == false) {
				return {
					success: result.success,
					message: result.message
				};
			}
		} catch (err) {
			console.log('Error creating todo:', err);
		}

		console.log('Create todo:', content);
	},
	deleteTodo: async ({ locals, request }) => {
		const formData = await request.formData();
		const todoId = formData.get('id');

		const token = locals.user?.token;

		try {
			const res = await fetch(`http://localhost:1234/api/v1/todos/${todoId}`, {
				method: 'DELETE',
				headers: {
					Authorization: `Bearer ${token}`
				}
			});
			const result = await res.json();
			console.log('Result:', result);
		} catch (err) {
			console.log('Error removing todo:', err);
		}

		console.log('Remove todo: ', todoId);
	}
};
