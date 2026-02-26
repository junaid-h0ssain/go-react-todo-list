import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { ChakraProvider } from '@chakra-ui/react'
import { system } from '@chakra-ui/react/preset'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ColorModeProvider } from './components/ui/color-mode'

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ChakraProvider value={system}>
      <ColorModeProvider>
        <QueryClientProvider client={queryClient}>
          <App />
        </QueryClientProvider>
      </ColorModeProvider>
    </ChakraProvider>
  </StrictMode>,
)
