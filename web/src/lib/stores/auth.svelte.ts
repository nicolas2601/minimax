import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { me } from '$lib/api/auth';
import type { User } from '$lib/schemas/auth';

const TOKEN_KEY = 'access_token';

interface AuthState {
	user: User | null;
	isLoading: boolean;
	isAuthenticated: boolean;
}

function createAuthStore() {
	let state = $state<AuthState>({
		user: null,
		isLoading: true,
		isAuthenticated: false
	});

	async function checkAuth(): Promise<void> {
		if (!browser) {
			state = { user: null, isLoading: false, isAuthenticated: false };
			return;
		}

		const token = localStorage.getItem(TOKEN_KEY);
		if (!token) {
			state = { user: null, isLoading: false, isAuthenticated: false };
			return;
		}

		try {
			const user = await me();
			state = { user, isLoading: false, isAuthenticated: true };
		} catch {
			// Token expired or invalid, try refresh via interceptor on next request
			state = { user: null, isLoading: false, isAuthenticated: false };
		}
	}

	function setUser(user: User): void {
		state = { user, isLoading: false, isAuthenticated: true };
	}

	function clearUser(): void {
		state = { user: null, isLoading: false, isAuthenticated: false };
	}

	function setLoading(loading: boolean): void {
		state = { ...state, isLoading: loading };
	}

	return {
		get user() { return state.user; },
		get isLoading() { return state.isLoading; },
		get isAuthenticated() { return state.isAuthenticated; },
		checkAuth,
		setUser,
		clearUser,
		setLoading
	};
}

export const authStore = createAuthStore();