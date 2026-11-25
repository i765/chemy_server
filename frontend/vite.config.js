import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
	plugins: [
		svelte({
			compilerOptions: {
				warningFilter: (warning) => {
					const ignored = [
						'a11y_click_events_have_key_events',
						'a11y_no_static_element_interactions'
					];
					return !ignored.includes(warning.code);
				}
			}
		})
	]
});
