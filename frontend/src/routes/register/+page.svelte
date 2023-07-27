<script>
	export let form;
</script>

<div>
  <h4>Register User</h4>
	<form method="POST">
		<label for="email">Email address</label>
		<input type="email" id="email" name="email" placeholder="Email address" value={form?.email ?? ''} required />

		<label for="password">Password</label>
		<input type="password" id="password" name="password" placeholder="Password" required />

		<label for="passwordConfirm">Password (confirm)</label>
		<input
			type="password"
			id="passwordConfirm"
			name="passwordConfirm"
			placeholder="Password"
			required
		/>

		<small
			>Your credentials are safe with us, we will never share your email with anyone else.</small
		>
		<button type="submit">Submit</button>
	</form>

	<div>
		{#if form?.success == false}
			<div class="error">{form?.message}</div>
		{/if}
	</div>
</div>

<style>
	.error {
		color: var(--del-color);
		font-weight: bold;
		text-align: center;
	}
</style>
