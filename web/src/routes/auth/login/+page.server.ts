import { LoginInputSchema } from '$lib/schemas/auth';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  // Return an empty default; the form is client-side handled via fetch
  return { defaultValues: { email: '', password: '' } };
};

// Use the schema at runtime on the client (no superValidate required).
export const _schema = LoginInputSchema;