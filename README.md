# **SuperCook**  

## Descripción
SuperCook es una aplicación web desarrollada en **Go** que implementa principios **RESTful** para gestionar alimentos, recetas y listas de compras. Utiliza **Gin Gonic** como framework para manejar solicitudes HTTP y **MongoDB** como base de datos para almacenar información.  

La interfaz de usuario está desarrollada en **HTML/CSS** y **JavaScript**, lo que permite una comunicación fluida entre el frontend y la lógica de negocio mediante solicitudes y respuestas en formato **JSON**. Este proyecto simula un entorno real de gestión de datos, poniendo en práctica conceptos clave del desarrollo web y la programación orientada a servicios.  

## Tecnologías utilizadas
- **Go**: Para la lógica de negocio y la implementación de controllers RESTful con Gin Gonic.  
- **MongoDB**: Como base de datos NoSQL para el almacenamiento y gestión de datos.  
- **HTML/CSS**: Para el diseño y la presentación de la interfaz de usuario.  
- **JavaScript**: Para manejar la comunicación entre el frontend y el backend.  
- **Docker**: Para el despliegue de la base de datos MongoDB en un contenedor.  

## Funcionalidades
### **Gestión de alimentos**:  
- Crear, editar y eliminar alimentos.  
- Definir propiedades como nombre, tipo, cantidad disponible, cantidad mínima, y precio por unidad.  

### **Recetas**:  
- Crear recetas con múltiples ingredientes y cantidades.  
- Buscar recetas disponibles según los alimentos existentes.  

### **Lista de compras**:  
- Generar listas de compras con los productos faltantes para completar una receta.  

### **Validaciones y Logs**:  
- Asegurar datos consistentes y registrar actividades para facilitar la depuración. 