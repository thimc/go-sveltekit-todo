<script>
	import { enhance } from '$app/forms';

	export let data;
	export let form;
</script>

<div>
	{#if data?.registeredEmail}
		<div class="alert">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				width="24"
				height="24"
				viewBox="0 0 24 24"
				style="margin-right: 10px"
				><path fill="var(--ins-color)" d="M9 16.17L4.83 12l-1.42 1.41L9 19L21 7l-1.41-1.41z" /></svg
			>
			<span style="color: var(--ins-color)"
				>Your account has been created! Please login to continue.</span
			>
		</div>
	{/if}
	<form method="POST">
		<div>
			<label for="email">Email</label>
			<input type="text" name="email" value={data?.registeredEmail ?? form?.email ?? ''} />
		</div>
		<div>
			<label for="password">Password</label>
			<input type="password" name="password" />
		</div>
		<button>Login</button>
		<small>No account? Click <a href="/register">here</a> to register one for free.</small>
	</form>
	<div>
		{#if form?.success == false}
			<div class="error">{form?.message}</div>
		{/if}
	</div>
</div>

<style>
	.alert {
		display: flex;
		border-radius: var(--border-radius);
		background: var(--code-background-color);
		padding: calc(var(--spacing) / 2) 0;
		justify-content: center;
		width: 100%;
		margin-top: 20px;
		padding-left: 20px;
		margin-bottom: 20px;
	}
	.error {
		color: var(--del-color);
		font-weight: bold;
		text-align: center;
	}
</style>
