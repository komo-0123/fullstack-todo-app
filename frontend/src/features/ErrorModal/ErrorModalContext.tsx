import { createContext } from "react";

// Context の型定義
export type ErrorModalContextType = {
  showError: (message: string) => void;
  closeError: () => void;
  errorMessage: string;
};

export const ErrorModalContext = createContext<ErrorModalContextType>({
  showError: () => {},
  closeError: () => {},
  errorMessage: "",
});
