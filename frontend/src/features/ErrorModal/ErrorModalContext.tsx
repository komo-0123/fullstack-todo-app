import { createContext } from "react";

// Context の型定義
export type ErrorModalContextType = {
  showError: (message: string) => void;
  closeError: () => void;
  errorMessage: string | null;
};

// Context を作成（このファイルは「Context の定義のみ」）
export const ErrorModalContext = createContext<ErrorModalContextType | undefined>(undefined);
