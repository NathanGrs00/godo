import React from "react";

/** A validation message under a form field. Renders nothing without a message. */
const FieldError: React.FC<{ id?: string; message?: string }> = ({ id, message }) =>
    message ? <p className="field-error" id={id} role="alert">{message}</p> : null;

export default FieldError;
