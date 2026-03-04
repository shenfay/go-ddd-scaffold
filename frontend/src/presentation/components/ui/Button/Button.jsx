import React from 'react';

const Button = ({ children, variant = 'primary', size = 'md', fullWidth = false, className = '', ...props }) => {
  const baseClasses = 'inline-flex items-center justify-center rounded transition-all focus:outline-none focus:ring-2 focus:ring-offset-2 font-sans';
  
  const variants = {
    primary: 'bg-primary hover:bg-primary-dark text-white focus:ring-primary',
    secondary: 'bg-secondary-dark hover:bg-secondary text-text-primary focus:ring-secondary',
    outline: 'border-1.5 border-primary bg-transparent hover:bg-primary hover:text-white text-primary focus:ring-primary',
    ghost: 'bg-transparent hover:bg-secondary text-text-primary focus:ring-secondary',
    danger: 'bg-danger hover:bg-danger-dark text-white focus:ring-danger',
  };
  
  const sizes = {
    sm: 'text-small py-2 px-3 rounded-sm',
    md: 'text-body py-2.5 px-5 rounded-md',
    lg: 'text-body py-3 px-6 rounded-lg',
  };
  
  const widthClass = fullWidth ? 'w-full' : '';
  const classes = `${baseClasses} ${variants[variant]} ${sizes[size]} ${widthClass} ${className}`;
  
  return (
    <button className={classes} {...props}>
      {children}
    </button>
  );
};

export default Button;
