import 'package:flutter/material.dart';
import '../widgets/app_drawer.dart';

class DashboardPage extends StatelessWidget {
  const DashboardPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const AppDrawer(),
      appBar: AppBar(title: const Text('Dashboard')),
      body: const Center(
        child: Text(
          'Dashboard Overview',
          style: TextStyle(fontSize: 18),
        ),
      ),
    );
  }
}
