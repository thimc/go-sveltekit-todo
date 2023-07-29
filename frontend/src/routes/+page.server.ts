import { API_URL } from '$env/static/private';

export const load = async ({ locals }) => {
	const { user } = locals;

	const getTodos = async (): Promise<App.Todo[] | null> => {
		console.log('Fetching todos');
		const token = user?.token ?? '';
		try {
			const res = await fetch(`${API_URL}/api/v1/todos`, {
				method: 'GET',
				headers: { Authorization: `Bearer ${token}` }
			})
				.then((response) => response.json())
				.then((data) => data.result as App.Todo[]);

			return res;
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

export const actions = {
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

		console.log('Remove todo:', todoId);
	},
	editTodo: async ({ locals, request }) => {
		const formData = await request.formData();
		const todoId = formData.get('id');
		const todoContent = formData.get('content');

		const token = locals.user?.token;
		try {
			const res = await fetch(`http://localhost:1234/api/v1/todos/${todoId}`, {
				method: 'PATCH',
				headers: {
					Authorization: `Bearer ${token}`
				},
				body: JSON.stringify({
					id: todoId,
					content: todoContent
				})
			});
			const result = await res.json();
			console.log('Result:', result);
		} catch (err) {
			console.log('Error removing todo:', err);
		}

		console.log('Edit todo:', todoId, todoContent);
	}
};
