import React from 'react';
import { Link } from 'react-router-dom'; // Assuming we'll use react-router for navigation

const HomePage = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col items-center justify-center p-4">
      <div className="max-w-md w-full bg-white rounded-xl shadow-md overflow-hidden md:max-w-2xl">
        <div className="p-8">
          <div className="flex items-center justify-center mb-6">
            <h1 className="text-3xl font-bold text-gray-800">Welcome to MathFun!</h1>
          </div>
          <p className="text-gray-600 text-center mb-8">
            Explore the world of mathematics through interactive 3D experiences.
          </p>
          <div className="flex justify-center">
            {/* Link to the examples page */}
            <Link
              to="/examples"
              className="bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-6 rounded-lg shadow-lg transition duration-300 ease-in-out transform hover:scale-105"
            >
              Start Exploring
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;