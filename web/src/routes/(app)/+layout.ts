import { redirect } from '@sveltejs/kit';
import { authStore } from '$lib/stores/auth.svelte.ts';

export async function load() {
	await authStore.checkAuth();

	if (!authStore.isAuthenticated) {
		throw redirect(302, '/auth/login');
	}

	return {};
}