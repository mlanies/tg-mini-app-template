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
    <div className="container" id="online-base-container">
      <div className="online-header2 bgrecolor">
        <div className="wrap-cont">
          <a className="back" href="/kiyanitsa/priceback?cid=uri2bchk8vdvk63eui1ip7vk5o"></a>
          <div className="title">
            <div>Выбор {selectedMaster ? 'услуг' : 'мастера'}</div>
            <small>{user ? `${user.first_name} ${user.last_name}` : ""}</small>
          </div>
          <div className="logo">
            <img src={telegramLogo} className="logo telegram" alt="Telegram logo" />
          </div>
        </div>
      </div>

      <div className="online-container">
        {/* Отображение мастеров или услуг */}
        {selectedMaster === null ? (
          <div className="online-masters">
            <div className="online-block online-services-wrap active">
              <div className="category bgrecolor">
                <div>Выберите мастера<span className="count"></span><span className="caret2"></span></div>
              </div>
              <div className="masters blkrecolor">
                {loading ? (
                  <p>Загрузка мастеров...</p>
                ) : (
                  masters.map((master) => (
                    <div
                      key={master.id}
                      className="master border p-4 rounded-lg shadow-md cursor-pointer hover:bg-gray-100 transition-all"
                      onClick={() => handleMasterSelect(master.id)}
                    >
                      <div className="data">
                        <div className="name font-bold text-lg">{master.name}</div>
                        <div className="info text-gray-600">
                          Опыт работы: {master.experience} лет
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        ) : (
          <div className="online-services">
            <div className="online-block online-services-wrap active">
              <div className="category bgrecolor">
                <div>Новые услуги<span className="count"></span><span className="caret2"></span></div>
              </div>
              <div className="services blkrecolor">
                {loading ? (
                  <p>Загрузка услуг...</p>
                ) : (
                  services.map((service) => (
                    <div
                      key={service.id}
                      className={`service ${selectedServices.includes(service.id) ? 'active' : ''} border p-4 rounded-lg shadow-md cursor-pointer hover:bg-gray-100 transition-all`}
                      data-id={service.id}
                      onClick={() => handleServiceSelect(service.id)}
                    >
                      <div className="data">
                        <div className="name">{service.name}</div>
                        <div className="info">
                          <span>от {service.price} ₽,</span>
                          <span>{service.duration} мин.</span>
                        </div>
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
          </div>
        )}

        {/* Кнопка для продолжения */}
        {selectedMaster !== null && (
          <div className="navbar-fixed-bottom" id="online-dialog-next">
            <div className="container vertical-align">
              <button
                type="button"
                className="text-white bg-gradient-to-br from-purple-600 to-blue-500 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-blue-300 dark:focus:ring-blue-800 font-medium rounded-lg text-sm px-5 py-2.5 text-center me-2 mb-2"
              >
                Purple to Blue
              </button>
              <button onClick={handleProceedBooking} className="online-def-btn bgrecolor">
                <span>Продолжить</span>
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
