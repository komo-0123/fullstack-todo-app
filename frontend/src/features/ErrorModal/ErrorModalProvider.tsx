import { useState } from "react";
import { ErrorModalContext } from "./ErrorModalContext";

export const ErrorModalProvider = ({ children }: { children: React.ReactNode }) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const showError = (message: string) => {
    setErrorMessage(message);
  };

  const closeError = () => {
    setErrorMessage(null);
  };

  return (
    <ErrorModalContext.Provider value={{ showError, closeError, errorMessage }}>
      {children}
    </ErrorModalContext.Provider>
  );
};
