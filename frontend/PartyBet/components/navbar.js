import React from "react";
import { View, Text, TouchableOpacity, StyleSheet } from "react-native";
import Icon from "react-native-vector-icons/MaterialIcons";

const Navbar = () => {
  return (
    <View style={styles.container}>
      <TouchableOpacity style={styles.button}>
        <Icon name="location-on" size={24} color="black" />
        <Text style={styles.label}>Past Events</Text>
      </TouchableOpacity>
      <TouchableOpacity style={styles.button}>
        <Icon name="add-circle-outline" size={24} color="black" />
        <Text style={styles.label}>Start New Event</Text>
      </TouchableOpacity>
      <TouchableOpacity style={styles.button}>
        <Icon name="notifications-none" size={24} color="black" />
        <Text style={styles.label}>Join Event</Text>
      </TouchableOpacity>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    flexDirection: "row",
    justifyContent: "space-around",
    alignItems: "center",
    borderTopWidth: 1,
    borderTopColor: "#dcdcdc",
    paddingVertical: 10,
    backgroundColor: "white", // Ensure the background is white to match the mockup
  },
  button: {
    alignItems: "center",
  },
  label: {
    marginTop: 4,
    fontSize: 12,
    color: "black",
  },
});

export default Navbar;
