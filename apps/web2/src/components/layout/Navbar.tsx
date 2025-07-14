'use client';

import { useState } from 'react';
import { useSession, signOut } from 'next-auth/react';
import { Menu, Transition } from '@headlessui/react';
import { Fragment } from 'react';
import {
  UserCircleIcon,
  QrCodeIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/react/24/outline';

export default function Navbar() {
  const { data: session } = useSession();
  const [showQR, setShowQR] = useState(false);

  const handleLogout = async () => {
    await signOut({ redirect: true, callbackUrl: '/login' });
  };

  return (
    <>
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="flex justify-between items-center h-16 px-6">
          <div className="flex items-center">
            <h2 className="text-xl font-semibold text-gray-800">管理面板</h2>
          </div>
          
          <div className="flex items-center space-x-4">
            <button
              onClick={() => setShowQR(true)}
              className="p-2 rounded-md hover:bg-gray-100 transition-colors"
              title="QR Code"
            >
              <QrCodeIcon className="h-5 w-5 text-gray-600" />
            </button>
            
            <Menu as="div" className="relative">
              <Menu.Button className="flex items-center space-x-2 p-2 rounded-md hover:bg-gray-100 transition-colors">
                <UserCircleIcon className="h-8 w-8 text-gray-600" />
                <span className="text-sm font-medium text-gray-700">
                  {session?.user?.name || session?.user?.email || '用戶'}
                </span>
              </Menu.Button>
              
              <Transition
                as={Fragment}
                enter="transition ease-out duration-100"
                enterFrom="transform opacity-0 scale-95"
                enterTo="transform opacity-100 scale-100"
                leave="transition ease-in duration-75"
                leaveFrom="transform opacity-100 scale-100"
                leaveTo="transform opacity-0 scale-95"
              >
                <Menu.Items className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                  <div className="py-1">
                    <Menu.Item>
                      {({ active }) => (
                        <button
                          className={`${
                            active ? 'bg-gray-100' : ''
                          } flex items-center px-4 py-2 text-sm text-gray-700 w-full text-left`}
                        >
                          <UserCircleIcon className="h-4 w-4 mr-2" />
                          修改個人資料
                        </button>
                      )}
                    </Menu.Item>
                    <Menu.Item>
                      {({ active }) => (
                        <button
                          onClick={handleLogout}
                          className={`${
                            active ? 'bg-gray-100' : ''
                          } flex items-center px-4 py-2 text-sm text-gray-700 w-full text-left`}
                        >
                          <ArrowRightOnRectangleIcon className="h-4 w-4 mr-2" />
                          登出
                        </button>
                      )}
                    </Menu.Item>
                  </div>
                </Menu.Items>
              </Transition>
            </Menu>
          </div>
        </div>
      </nav>
      
      {showQR && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded-lg shadow-xl">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">手機掃描登入</h3>
              <button
                onClick={() => setShowQR(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                ✕
              </button>
            </div>
            <div className="flex items-center justify-center w-64 h-64 bg-gray-100 rounded-lg">
              <span className="text-gray-500">QR Code 將在此顯示</span>
            </div>
          </div>
        </div>
      )}
    </>
  );
}