"use client";

import {
  environmentManager,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

function makeQueryCLient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        //above 0 to avoid refetching immediately on the client
        staleTime: 60 * 1000,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined = undefined;

function getQueryClient() {
  if (environmentManager.isServer()) {
    return makeQueryCLient();
  } else {
    if (!browserQueryClient) browserQueryClient = makeQueryCLient();
    return browserQueryClient;
  }
}

const queryClient = getQueryClient();

export default function TanstackProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <QueryClientProvider client={queryClient}>
      {children} <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}
