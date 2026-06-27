import { RegisterInputSchema } from '$lib/schemas/auth';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  return { defaultValues: { email: '', password: '', display_name: '' } };
};

export const _schema = RegisterInputSchema;