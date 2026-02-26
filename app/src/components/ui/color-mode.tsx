import { ThemeProvider, useTheme } from "next-themes";
import type { ThemeProviderProps } from "next-themes";

export function ColorModeProvider(props: ThemeProviderProps) {
	return <ThemeProvider attribute='class' disableTransitionOnChange {...props} />;
}

export function useColorMode() {
	const { resolvedTheme, setTheme } = useTheme();
	const colorMode = resolvedTheme === "dark" ? "dark" : "light";

	const toggleColorMode = () => {
		setTheme(colorMode === "dark" ? "light" : "dark");
	};

	return {
		colorMode,
		setColorMode: setTheme,
		toggleColorMode,
	};
}

export function useColorModeValue<T>(light: T, dark: T) {
	const { colorMode } = useColorMode();
	return colorMode === "dark" ? dark : light;
}
