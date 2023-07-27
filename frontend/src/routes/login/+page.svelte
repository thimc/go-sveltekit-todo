<script>
	import { enhance } from '$app/forms';

	export let form;
</script>

<div>
	<form method="POST" action="?/login" use:enhance>
		<div>
			<label for="email">Email</label><input
				type="text"
				name="email"
				value={form?.email ?? 'user@domain.com'}
			/>
		</div>
		<div>
			<label for="password">Password</label>
			<input type="password" name="password" value="string" />
		</div>
		<button>Login</button>
	</form>
	{#if form?.success == false}
		<div class="error">Error: {form?.message}</div>
	{/if}

	<pre>{JSON.stringify(form?.cookie, null, 2)}</pre>
</div>

<style>
	.error {
		color: var(--del-color);
		font-weight: bold;
	}
</style>
