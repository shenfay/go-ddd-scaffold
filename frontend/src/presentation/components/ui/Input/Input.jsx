import React from 'react';

const Input = ({ label, className = '', ...props }) => {
  const inputClasses = `input ${className}`;
  
  return (
    <div className="mb-4">
      {label && (
        <label className="block text-body font-semibold mb-2 text-text-primary">
          {label}
        </label>
      )}
      <input className={inputClasses} {...props} />
    </div>
  );
};

export default Input;