import { useEffect, useState } from 'react';
import telegramLogo from './assets/telegram.svg';
import './output.css';
import './bootstrap.min.css';

function App() {
  const [masters, setMasters] = useState([]);
  const [services, setServices] = useState([]);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [selectedServices, setSelectedServices] = useState([]);
  const [selectedMaster, setSelectedMaster] = useState(null);

  useEffect(() => {
    // Настройка темы Telegram
    const themeParams = window.Telegram.WebApp.themeParams;
    document.body.style.backgroundColor = themeParams.bg_color || '#ffffff';
    document.body.style.color = themeParams.text_color || '#000000';

    // Извлечение данных пользователя из WebApp
    const initData = window.Telegram.WebApp.initDataUnsafe;
    if (initData && initData.user) {
      setUser(initData.user);
      console.log("Данные пользователя:", initData.user);
    } else {
      console.error("Данные пользователя не переданы.");
    }

    // Загрузка списка мастеров из API
    const fetchMasters = async () => {
      setLoading(true);
      try {
        const response = await fetch('http://localhost:3000/api/masters'); // Вызов вашего API для мастеров
        if (!response.ok) {
          throw new Error(`Ошибка загрузки мастеров: ${response.statusText}`);
        }
        const data = await response.json();
        setMasters(data);
      } catch (error) {
        console.error("Ошибка загрузки мастеров:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchMasters();
  }, []);

  const fetchServices = async (masterId) => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:3000/api/masters/${masterId}/services`); // Вызов вашего API для получения услуг выбранного мастера
      if (!response.ok) {
        throw new Error(`Ошибка загрузки услуг: ${response.statusText}`);
      }
      const data = await response.json();
      setServices(data);
    } catch (error) {
      console.error("Ошибка загрузки услуг:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleMasterSelect = (masterId) => {
    setSelectedMaster(masterId);
    fetchServices(masterId); // Загрузите услуги для выбранного мастера
  };

  const handleServiceSelect = (serviceId) => {
    setSelectedServices((prevSelected) => {
      if (prevSelected.includes(serviceId)) {
        return prevSelected.filter((id) => id !== serviceId);
      } else {
        return [...prevSelected, serviceId];
      }
    });
  };

  const handleProceedBooking = () => {
    if (selectedServices.length === 0) {
      window.Telegram.WebApp.showAlert("Пожалуйста, выберите хотя бы одну услугу для продолжения.");
      return;
    }

    const bookingDetails = selectedServices.map((serviceId) => {
      const service = services.find((s) => s.id === serviceId);
      return service ? service.name : "";
    }).join(', ');

    window.Telegram.WebApp.showAlert(`Вы выбрали следующие услуги: ${bookingDetails}`);
  };

  return (
    <div className="container mx-auto py-8" id="online-base-container">
      <div className="bg-gradient-to-r from-blue-600 to-purple-600 text-white p-6 rounded-lg shadow-lg mb-6">
        <div className="flex items-center justify-between">
          <a className="text-white font-bold text-lg hover:underline" href="/kiyanitsa/priceback?cid=uri2bchk8vdvk63eui1ip7vk5o">
            ← Назад
          </a>
          <div className="text-center">
            <h1 className="text-2xl font-semibold">{selectedMaster ? 'Выбор услуг' : 'Выбор мастера'}</h1>
            <small>{user ? `${user.first_name} ${user.last_name}` : ""}</small>
          </div>
          <div className="flex-shrink-0">
            <img src={telegramLogo} className="w-12 h-12 rounded-full shadow-md" alt="Telegram logo" />
          </div>
        </div>
      </div>

      <div className="online-container">
        {/* Отображение мастеров или услуг */}
        {selectedMaster === null ? (
          <div className="online-masters">
            <div className="bg-white p-6 rounded-lg shadow-md mb-6">
              <h2 className="text-lg font-bold mb-4">Выберите мастера</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {loading ? (
                  <p>Загрузка мастеров...</p>
                ) : (
                  masters.map((master) => (
                    <div
                      key={master.id}
                      className="border p-4 rounded-lg shadow-md cursor-pointer hover:bg-blue-50 transition-all"
                      onClick={() => handleMasterSelect(master.id)}
                    >
                      <h3 className="font-bold text-xl mb-2">{master.name}</h3>
                      <p className="text-gray-600">Опыт работы: {master.experience} лет</p>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        ) : (
          <div className="online-services">
            <div className="bg-white p-6 rounded-lg shadow-md mb-6">
              <h2 className="text-lg font-bold mb-4">Выберите услуги</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {loading ? (
                  <p>Загрузка услуг...</p>
                ) : (
                  services.map((service) => (
                    <div
                      key={service.id}
                      className={`border p-4 rounded-lg shadow-md cursor-pointer hover:bg-blue-50 transition-all ${selectedServices.includes(service.id) ? 'bg-blue-100 border-blue-400' : 'border-gray-200'}`}
                      onClick={() => handleServiceSelect(service.id)}
                    >
                      <h3 className="font-bold text-xl mb-2">{service.name}</h3>
                      <p className="text-gray-600">от {service.price} ₽, {service.duration} мин.</p>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        )}

        {/* Кнопка для продолжения */}
        {selectedMaster !== null && (
          <div className="fixed bottom-0 left-0 w-full bg-white shadow-lg p-4">
            <div className="flex justify-center">
              <button
                onClick={handleProceedBooking}
                className="bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition duration-200"
              >
                Продолжить
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
