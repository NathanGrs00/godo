import React from "react";
import { useEscapeKey } from "../../hooks/useEscapeKey.ts";
import "../../styles/Common/Modal.css";

interface ModalProps {
    onClose: () => void;
    children: React.ReactNode;
}

/** Shows children above a dimmed backdrop. Clicking the backdrop or pressing Escape closes it. */
const Modal: React.FC<ModalProps> = ({ onClose, children }) => {
    useEscapeKey(onClose);

    return (
        <>
            <div className="modal-backdrop" data-testid="modal-backdrop" onClick={onClose} />
            {children}
        </>
    );
};

export default Modal;
