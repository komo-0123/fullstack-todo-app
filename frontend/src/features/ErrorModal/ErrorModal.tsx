import { useEffect, useRef, useContext } from "react";
import { createPortal } from "react-dom";
import { ErrorModalContext } from "./ErrorModalContext";

const ErrorModal = () => {
  const { errorMessage, closeError } = useContext(ErrorModalContext);
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    if (errorMessage) {
      dialogRef.current?.showModal();
    }
  }, [errorMessage]);

  if (!errorMessage) {
    return null;
  }

  return createPortal(
    <dialog
      ref={dialogRef}
      className="backdrop:bg-black/50 fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-white p-6 rounded-lg shadow-lg"
    >
      <p className="text-red-600 text-lg">{errorMessage}</p>
      <button
        onClick={() => {
          dialogRef.current?.close();
          closeError();
        }}
        className="mt-4 px-4 py-2 bg-teal-500 text-white rounded hover:bg-teal-600 mx-auto block"
      >
        OK
      </button>
    </dialog>,
    document.body
  );
};

export default ErrorModal;
