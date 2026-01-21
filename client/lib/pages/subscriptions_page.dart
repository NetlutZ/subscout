import 'package:client/pages/add_subscription_page.dart';
import 'package:client/widgets/app_drawer.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../models/subscription.dart';
import '../services/subscription_service.dart';

class SubscriptionsPage extends StatefulWidget {
  const SubscriptionsPage({super.key});

  @override
  State<SubscriptionsPage> createState() => _SubscriptionsPageState();
}

class _SubscriptionsPageState extends State<SubscriptionsPage> {
  final _service = SubscriptionService();
  late Future<List<Subscription>> _future;

  @override
  void initState() {
    super.initState();
    _loadSubscriptions();
  }

  void _loadSubscriptions() {
    _future = _service.fetchSubscriptions();
  }

  Future<void> _deleteSubscription(int id) async {
    await _service.deleteSubscription(id);
    setState(() {
      _loadSubscriptions(); // refresh list
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const AppDrawer(),
      appBar: AppBar(title: const Text('Subscriptions')),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          final result = await Navigator.push(
            context,
            MaterialPageRoute(builder: (_) => const AddSubscriptionPage()),
          );

          if (result == true) {
            setState(() => _loadSubscriptions());
          }
        },
        child: const Icon(Icons.add),
      ),
      body: FutureBuilder<List<Subscription>>(
        future: _future,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }

          if (snapshot.hasError) {
            return Center(child: Text(snapshot.error.toString()));
          }

          final subscriptions = snapshot.data!;

          if (subscriptions.isEmpty) {
            return const Center(child: Text('No subscriptions found'));
          }

          return ListView.builder(
            itemCount: subscriptions.length,
            itemBuilder: (_, index) {
              final sub = subscriptions[index];

              return Card(
                margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                child: ListTile(
                  title: Text(sub.name),
                  subtitle: Text(
                    '${sub.category} • ${sub.billingCycle}\nNext: ${DateFormat('dd/MM/yyyy').format(sub.billingDate)}',
                  ),
                  trailing: SizedBox(
                    width: 120, // constrain width
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        Text(
                          '${sub.currency} ${sub.amount}',
                          style: const TextStyle(fontWeight: FontWeight.bold),
                        ),
                        IconButton(
                          icon: const Icon(Icons.delete, color: Colors.red),
                          onPressed: () => _deleteSubscription(sub.id),
                        ),
                      ],
                    ),
                  ),
                  isThreeLine: true,
                ),
              );
            },
          );
        },
      ),
    );
  }
}